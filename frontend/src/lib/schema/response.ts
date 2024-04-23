import { z } from "zod";

export const successSchema = z.object({
  message: z.string(),
  data: z.any().optional(),
});

export const errorSchema = z.object({
  message: z.string(),
});

export const errorValidationSchema = z.object({
  message: z.string(),
  errors: z.array(
    z.object({
      field: z.string(),
      error: z.string(),
    })
  ),
});

export type Success = z.infer<typeof successSchema>;
export type Error = z.infer<typeof errorSchema>;
export type ErrorValidation = z.infer<typeof errorValidationSchema>;
