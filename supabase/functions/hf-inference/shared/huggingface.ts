// import { ProcessingResult, MODEL_CONFIG } from './types.ts';
//
// export class HuggingFaceClient {
//     private readonly apiKey: string;
//     private readonly headers: HeadersInit;
//
//     constructor(apiKey: string) {
//         if (!apiKey) {
//             throw new Error('HuggingFace API token is required');
//         }
//
//         this.apiKey = apiKey;
//         this.headers = {
//             'Authorization': `Bearer ${apiKey}`,
//             'Content-Type': 'application/json'
//         };
//     }
//
//     async processText(input: string): Promise<ProcessingResult> {
//         try {
//             const config = MODEL_CONFIG.text2text;
//             const response = await fetch(config.endpoint, {
//                 method: 'POST',
//                 headers: this.headers,
//                 body: JSON.stringify({
//                     inputs: input,
//                     parameters: config.parameters
//                 })
//             });
//
//             if (!response.ok) {
//                 throw new Error(`HuggingFace API error: ${response.statusText}`);
//             }
//
//             const result = await response.json();
//             return {
//                 status: 'success',
//                 data: {
//                     output: result.generated_text,
//                     processingTime: Date.now(),
//                     tokenCount: result.usage?.total_tokens
//                 }
//             };
//         } catch (error) {
//             return {
//                 status: 'error',
//                 error: error as string
//             };
//         }
//     }
//
//     async processImage(input: string): Promise<ProcessingResult> {
//         try {
//             const config = MODEL_CONFIG.text2image;
//             const response = await fetch(config.endpoint, {
//                 method: 'POST',
//                 headers: this.headers,
//                 body: JSON.stringify({
//                     inputs: input,
//                     parameters: config.parameters
//                 })
//             });
//
//             if (!response.ok) {
//                 throw new Error(`HuggingFace API error: ${response.statusText}`);
//             }
//
//             // Image response comes as binary data
//             const imageData = new Uint8Array(await response.arrayBuffer());
//
//             return {
//                 status: 'success',
//                 data: {
//                     output: imageData,
//                     processingTime: Date.now()
//                 }
//             };
//         } catch (error) {
//             return {
//                 status: 'error',
//                 error: error as string
//             };
//         }
//     }
// }
//
//     async processImage(input: string): Promise<ProcessingResult> {
//         try {
//             // Validate input
//             if (!input || typeof input !== 'string') {
//                 throw new Error('Invalid input: Input must be a non-empty string');
//             }
//
//             // Get model configuration
//             const config = MODEL_CONFIG.text2image;
//             if (!config || !config.endpoint) {
//                 throw new Error('Invalid model configuration: Missing endpoint');
//             }
//
//             // Log request details (helpful for debugging)
//             console.log('Making request to HuggingFace API:', {
//                 endpoint: config.endpoint,
//                 input: input,
//                 parameters: config.parameters
//             });
//
//             // Make the API request
//             const response = await fetch(config.endpoint, {
//                 method: 'POST',
//                 headers: this.headers,
//                 body: JSON.stringify({
//                     inputs: input,
//                     parameters: {
//                         ...config.parameters,
//                         return_full_text: false
//                     }
//                 })
//             });
//
//             // Check for HTTP errors
//             if (!response.ok) {
//                 const errorText = await response.text();
//                 throw new Error(
//                     `HuggingFace API error (${response.status}): ${response.statusText}\n${errorText}`
//                 );
//             }
//
//             // Verify content type is image
//             const contentType = response.headers.get('content-type');
//             if (!contentType?.includes('image/')) {
//                 throw new Error(
//                     `Unexpected content type: ${contentType}. Expected image response.`
//                 );
//             }
//
//             // Process image response
//             const imageData = new Uint8Array(await response.arrayBuffer());
//             if (!imageData.length) {
//                 throw new Error('Received empty image data from API');
//             }
//
//             return {
//                 status: 'success',
//                 data: {
//                     output: imageData,
//                     processingTime: Date.now(),
//                     //contentType: contentType
//                 }
//             };
//         } catch (error) {
//             // Provide detailed error information
//             console.error('Image processing error:', error);
//             return {
//                 status: 'error',
//                 error: error instanceof Error ?
//                     {
//                         message: error.message,
//                         name: error.name,
//                         stack: error.stack
//                     } :
//                     String(error)
//             };
//         }
//     }
//
// }

export class HuggingFaceClient {
    private readonly headers: HeadersInit;
    private readonly endpoints = {
        text2text: 'https://api-inference.huggingface.co/models/google/gemma-2-2b-it',
        text2image: 'https://api-inference.huggingface.co/models/black-forest-labs/FLUX.1-dev'
    };

    constructor(apiKey: string) {
        this.headers = {
            'Authorization': `Bearer ${apiKey}`,
            'Content-Type': 'application/json'
        };
    }

    async processText(input: string) {
        const response = await fetch(this.endpoints.text2text, {
            method: 'POST',
            headers: this.headers,
            body: JSON.stringify({
                inputs: input,
                parameters: {
                    max_new_tokens: 512,
                    temperature: 0.7,
                    top_p: 0.95,
                    return_full_text: false
                }
            })
        });

        if (!response.ok) {
            throw new Error('Text generation failed');
        }

        return await response.json();
    }

    async processImage(input: string) {
        const response = await fetch(this.endpoints.text2image, {
            method: 'POST',
            headers: this.headers,
            body: JSON.stringify({
                inputs: input,
                parameters: {
                    num_inference_steps: 30,
                    guidance_scale: 7.5
                }
            })
        });

        if (!response.ok) {
            throw new Error('Image generation failed');
        }

        return {
            data: new Uint8Array(await response.arrayBuffer())
        };
    }
}