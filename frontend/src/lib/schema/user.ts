import { z } from "zod";
import { metadataSchema } from "./metadata";

export const userSchema = z.object({
  id: z.number(),
  name: z.string(),
  email: z.string(),
  password: z.string(),
  bio: z.string().optional(),
  avatar_url: z
    .string()
    .url()
    .endsWith(".jpg" || ".png")
    .optional()
    .or(z.literal("")),
  created_at: z.string().datetime(),
  updated_at: z.string().datetime(),
  tokens: z
    .object({
      access_token: z.string().optional(),
      refresh_token: z.string().optional(),
    })
    .optional(),
});

export const userRequestSchema = z.object({
  name: z.string().min(4).max(20).optional(),
  email: z.string().optional(),
  password: z.string().optional(),
  bio: z.string().max(50).optional(),
  avatar_url: z
    .string()
    .url()
    .endsWith(".jpg")
    .or(z.string().endsWith(".png"))
    .optional()
    .or(z.literal("")),
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
export type UserRequest = z.infer<typeof userRequestSchema>;
export type UserPagination = z.infer<typeof userPaginationSchema>;
