import { createClient } from "jsr:@supabase/supabase-js";
import { HuggingFaceClient } from "./huggingface.ts";
import { ModelRequest, ProcessingResult } from "./types.ts";

export class RequestProcessor {
    private supabase: any;
    private hfClient: HuggingFaceClient;

    constructor(supabaseUrl: string, supabaseKey: string, hfApiKey: string) {
        const serviceRoleKey = Deno.env.get("SUPABASE_SERVICE_ROLE_KEY");
        const keyToUse = serviceRoleKey ? serviceRoleKey : supabaseKey;
        console.log("Using key for supabase client:", serviceRoleKey ? "service_role" : "anon");
        this.supabase = createClient(supabaseUrl, keyToUse);
        this.hfClient = new HuggingFaceClient(hfApiKey);
    }

    private async updateRequestStatus(
        requestId: string,
        status: string,
        data?: any,
        error?: string
    ) {
        const updates = {
            status,
            ...(data && { output_data: data }),
            ...(error && { error_msg: error }),
            ...(data?.processingTime && { processing_time: data.processingTime }),
            ...(data?.tokenCount && { token_count: data.tokenCount })
        };
        await this.supabase
            .from("model_requests")
            .update(updates)
            .match({ id: requestId });
    }

    async processRequest(request: ModelRequest): Promise<ProcessingResult> {
        try {
            await this.updateRequestStatus(request.id, "processing");
            const result = await this.hfClient.processText(request.input);
            if (result.status === "error") {
                await this.updateRequestStatus(request.id, "failed", null, result.error);
                return result;
            }
            await this.updateRequestStatus(request.id, "completed", result.data);
            return result;
        } catch (error) {
            await this.updateRequestStatus(
                request.id,
                "failed",
                null,
                error instanceof Error ? error.message : String(error)
            );
            return {
                status: "error",
                error: error instanceof Error ? error.message : String(error)
            };
        }
    }
}
