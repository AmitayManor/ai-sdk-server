export interface ModelRequest {
    id: string;
    input: string;
    modelType: "text2text";
    userId: string;
}

export interface ModelResponse {
    output: string;
    processingTime: number;
    tokenCount?: number;
}

export interface ProcessingResult {
    status: "success" | "error";
    data?: ModelResponse;
    error?: string;
}

export const MODEL_CONFIG = {
    text2text: {
        endpoint: "https://api-inference.huggingface.co/models/google/gemma-2-2b-it",
        params: { max_new_tokens: 512, temperature: 0.7, top_p: 0.95, return_full_text: false }
    }
} as const;
