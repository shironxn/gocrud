import { z } from "zod";

export const metadataSchema = z.object({
  sort: z.enum(["id", "user_id", "name", "title", "created_at", "updated_at"]),
  order: z.enum(["asc", "desc"]),
  total_records: z.number(),
  total_pages: z.number(),
  limit: z.number(),
  page: z.number(),
});

export type Metadata = z.infer<typeof metadataSchema>;
