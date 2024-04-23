import { z } from "zod";

export const authRegisterSchema = z.object({
  name: z.string().min(4).max(30),
  email: z.string().email(),
  password: z.string().min(8).max(100),
});

export const authLoginSchema = z.object({
  email: z.string().email(),
  password: z.string().min(8),
});

export type AuthLogin = z.infer<typeof authLoginSchema>;
export type AuthRegister = z.infer<typeof authRegisterSchema>;
