import { z } from "zod";
import { metadataSchema } from "./metadata";

export const noteSchema = z.object({
  id: z.number(),
  title: z.string().max(30),
  description: z.string().max(100),
  cover_url: z
    .string()
    .url()
    .endsWith(".jpg" || ".png"),
  content: z.string(),
  visibility: z.enum(["public", "private"]),
  author: z.object({
    id: z.number(),
    name: z.string(),
    avatar_url: z
      .string()
      .url()
      .endsWith(".jpg" || ".png"),
  }),
  created_at: z.string().datetime(),
  updated_at: z.string().datetime(),
});

export const noteRequestSchema = z.object({
  title: z.string().min(1).max(25),
  description: z.string().min(1).max(50),
  cover_url: z
    .string()
    .url()
    .endsWith(".jpg" || ".png"),
  content: z.string().min(1),
  visibility: z.enum(["public", "private"]),
});

export const noteQuerySchema = z.object({
  title: z.string(),
  visibility: z.enum(["public", "private"]),
  user_id: z.number(),
});

export const notePaginationSchema = z.object({
  notes: z.array(noteSchema),
  metadata: metadataSchema,
});

export type Note = z.infer<typeof noteSchema>;
export type NoteRequest = z.infer<typeof noteRequestSchema>;
export type NoteQuery = z.infer<typeof noteQuerySchema>;
export type NotePagination = z.infer<typeof notePaginationSchema>;
