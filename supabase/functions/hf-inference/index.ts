// import { serve } from 'https://deno.land/std@0.168.0/http/server.ts';
// import { RequestProcessor } from './shared/processor.ts';
// import { ModelRequest } from './shared/types.ts';
//
// const processor = new RequestProcessor(
//     Deno.env.get('_SUPABASE_URL')!,
//     Deno.env.get('_SUPABASE_ANON_KEY')!,
//     Deno.env.get('HF_API_TOKEN')!
// );
//
// serve(async (req: Request) => {
//   try {
//     if (req.method !== 'POST') {
//       return new Response(
//           JSON.stringify({ error: 'Method not allowed' }),
//           { status: 405, headers: { 'Content-Type': 'application/json' } }
//       );
//     }
//
//     const request = await req.json() as ModelRequest;
//     if (!request.id || !request.modelType || !request.input || !request.userId) {
//       return new Response(
//           JSON.stringify({ error: 'Invalid request format' }),
//           { status: 400, headers: { 'Content-Type': 'application/json' } }
//       );
//     }
//
//     const result = await processor.processRequest(request);
//
//     return new Response(
//         JSON.stringify(result),
//         {
//           status: result.status === 'success' ? 200 : 500,
//           headers: { 'Content-Type': 'application/json' }
//         }
//     );
//
//   } catch (error) {
//     return new Response(
//         JSON.stringify({ error: error}),
//         { status: 500, headers: { 'Content-Type': 'application/json' } }
//     );
//   }
// });
//
// import { serve } from 'https://deno.land/std@0.168.0/http/server.ts';
// import { RequestProcessor } from './shared/processor.ts';
// import { ModelRequest } from './shared/types.ts';
//
// // Let's add better error handling and logging
// const processor = new RequestProcessor(
//     // These environment variables must be set in your Supabase dashboard
//     Deno.env.get('SUPABASE_URL') ?? '',
//     Deno.env.get('SUPABASE_ANON_KEY') ?? '',
//     Deno.env.get('HF_API_TOKEN') ?? ''
// );
//
// serve(async (req: Request) => {
//   try {
//     // Only allow POST requests
//     if (req.method !== 'POST') {
//       return new Response(
//           JSON.stringify({ error: 'Method not allowed' }),
//           { status: 405, headers: { 'Content-Type': 'application/json' } }
//       );
//     }
//
//     // Log the incoming request
//     console.log('Received request:', {
//       method: req.method,
//       headers: Object.fromEntries(req.headers.entries())
//     });
//
//     // Parse the request body
//     let request: ModelRequest;
//     try {
//       request = await req.json();
//       console.log('Parsed request body:', request);
//     } catch (e) {
//       console.error('Failed to parse request body:', e);
//       return new Response(
//           JSON.stringify({ error: 'Invalid JSON in request body' }),
//           { status: 400, headers: { 'Content-Type': 'application/json' } }
//       );
//     }
//
//     // Validate required fields
//     if (!request.id || !request.input || !request.modelType || !request.userId) {
//       const missingFields = [];
//       if (!request.id) missingFields.push('id');
//       if (!request.input) missingFields.push('input');
//       if (!request.modelType) missingFields.push('modelType');
//       if (!request.userId) missingFields.push('userId');
//
//       console.error('Missing required fields:', missingFields);
//       return new Response(
//           JSON.stringify({
//             error: `Missing required fields: ${missingFields.join(', ')}`
//           }),
//           { status: 400, headers: { 'Content-Type': 'application/json' } }
//       );
//     }
//
//     // Process the request
//     console.log('Starting request processing...');
//     const result = await processor.processRequest(request);
//     console.log('Processing completed:', result);
//
//     return new Response(
//         JSON.stringify(result),
//         {
//           status: result.status === 'success' ? 200 : 500,
//           headers: { 'Content-Type': 'application/json' }
//         }
//     );
//
//   } catch (error) {
//     // Log the full error details
//     console.error('Unhandled error in edge function:', {
//       error: error,
//       message: error,
//       stack: error
//     });
//
//     return new Response(
//         JSON.stringify({
//           error: 'Internal server error',
//           details: error,
//           // Include stack trace in development, remove in production
//           stack: Deno.env.get('ENVIRONMENT') === 'development' ? error : undefined
//         }),
//         {
//           status: 500,
//           headers: { 'Content-Type': 'application/json' }
//         }
//     );
//   }
// });

import { serve } from 'https://deno.land/std@0.168.0/http/server.ts';
import { createClient } from 'jsr:@supabase/supabase-js';
import { HuggingFaceClient } from './shared/huggingface.ts';
import { ModelRequest } from './shared/types.ts';

serve(async (req: Request) => {
  try {
    // Initialize clients
    const supabase = createClient(
        Deno.env.get('SUPABASE_URL') ?? '',
        Deno.env.get('SUPABASE_ANON_KEY') ?? ''
    );

    const hf = new HuggingFaceClient(Deno.env.get('HF_API_TOKEN') ?? '');

    // Parse request
    const { id, input, modelType, userId } = await req.json() as ModelRequest;

    // Update status to processing
    await supabase.from('model_requests')
        .update({ status: 'processing' })
        .eq('id', id);

    // Process based on model type
    const result = modelType === 'text2text'
        ? await hf.processText(input)
        : await hf.processImage(input);

    // For image results, store in Supabase Storage
    let outputData = result;
    if (modelType === 'text2image' && result.data) {
      const filePath = `${userId}/${id}.png`;
      await supabase.storage
          .from('generated-images')
          .upload(filePath, result.data, { contentType: 'image/png' });
      outputData = { path: filePath };
    }

    // Update request with results
    await supabase.from('model_requests')
        .update({
          status: 'completed',
          output_data: outputData
        })
        .eq('id', id);

    return new Response(JSON.stringify(outputData), {
      headers: { 'Content-Type': 'application/json' }
    });

  } catch (error) {
    return new Response(JSON.stringify({ error: error}), {
      status: 500,
      headers: { 'Content-Type': 'application/json' }
    });
  }
});