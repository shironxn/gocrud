import { z } from "zod";
import { metadataSchema } from "./metadata";

export const userSchema = z.object({
  id: z.number(),
  name: z.string(),
  email: z.string(),
  bio: z.string().optional(),
  avatar_url: z
    .string()
    .url()
    .endsWith(".jpg" || ".png")
    .optional(),
  created_at: z.string().datetime(),
  updated_at: z.string().datetime(),
  tokens: z
    .object({
      access_token: z.string().optional(),
      refresh_token: z.string().optional(),
    })
    .optional(),
});

export const userTokenSchema = z.object({
  access_token: z.string(),
  refresh_token: z.string(),
});

export const claimsSchema = z.object({
  user_id: z.number(),
  exp: z.string().datetime(),
});

export const userPaginationSchema = z.object({
  users: z.array(userSchema),
  metadata: metadataSchema,
});

export type UserToken = z.infer<typeof userTokenSchema>;
export type Claims = z.infer<typeof claimsSchema>;
export type User = z.infer<typeof userSchema>;
export type UserPagination = z.infer<typeof userPaginationSchema>;
