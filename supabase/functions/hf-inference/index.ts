import { serve } from 'https://deno.land/std@0.168.0/http/server.ts';
import { RequestProcessor } from './shared/processor.ts';
import { ModelRequest } from './shared/types.ts';

const supabaseUrl = Deno.env.get('SUPABASE_URL') || '';
const supabaseKey = Deno.env.get('SUPABASE_ANON_KEY')|| '';
const hfApiKey = Deno.env.get('HF_API_TOKEN')|| '';

const processor = new RequestProcessor(supabaseUrl, supabaseKey, hfApiKey);

console.log("Edge function started and waiting for requests...");

serve(async (req: Request) => {
  console.log(`Received ${req.method} request`);
  try {
    if (req.method !== 'POST') {
      return new Response(
          JSON.stringify({ error: 'Method not allowed' }),
          { status: 405, headers: { 'Content-Type': 'application/json' } }
      );
    }

    const request = await req.json() as ModelRequest;
    console.log("Request payload:", request);

    if (!request.id || !request.modelType || !request.input || !request.userId) {
      console.error("Invalid request payload", request);
      return new Response(
          JSON.stringify({ error: 'Invalid request format' }),
          { status: 400, headers: { 'Content-Type': 'application/json' } }
      );
    }

    const result = await processor.processRequest(request);
    console.log("Processing result:", result);

    return new Response(
        JSON.stringify(result),
        {
          status: result.status === 'success' ? 200 : 500,
          headers: { 'Content-Type': 'application/json' }
        }
    );

  } catch (error) {
    console.error("Error processing request:", error);
    return new Response(
        JSON.stringify({ error: error instanceof Error ? error.message : error }),
        { status: 500, headers: { 'Content-Type': 'application/json' } }
    );
  }
});
