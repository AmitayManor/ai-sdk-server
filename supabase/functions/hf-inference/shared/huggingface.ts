import { MODEL_CONFIG } from "./types.ts";
import { ProcessingResult } from "./types.ts";

export class HuggingFaceClient {
    private readonly apiKey: string;
    private readonly headers: HeadersInit;

    constructor(apiKey: string) {
        this.apiKey = apiKey;
        this.headers = {
            "Authorization": `Bearer ${apiKey}`,
            "Content-Type": "application/json",
            "Accept": "application/json"
        };
    }

    async processText(input: string): Promise<ProcessingResult> {
        try {
            const config = MODEL_CONFIG.text2text;
            const response = await fetch(config.endpoint, {
                method: "POST",
                headers: this.headers,
                body: JSON.stringify({
                    inputs: input,
                    parameters: config.params
                })
            });
            if (!response.ok) {
                throw new Error(`HuggingFace API error: ${response.statusText}`);
            }
            const result = await response.json();
            console.log("Text generation response:", result);
            let generated = "";
            if (Array.isArray(result) && result.length > 0 && result[0].generated_text) {
                generated = result[0].generated_text;
            } else if (result.generated_text) {
                generated = result.generated_text;
            } else {
                throw new Error("No generated_text found in response");
            }
            return {
                status: "success",
                data: {
                    output: generated,
                    processingTime: Date.now()
                }
            };
        } catch (error) {
            return {
                status: "error",
                error: error instanceof Error ? error.message : String(error)
            };
        }
    }
}
