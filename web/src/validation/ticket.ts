import { z } from "zod";

export const MAX_FILE_SIZE = 2 * 1024 * 1024;
export const ACCEPTED_EXTENSIONS = ["jpg", "jpeg", "png"];

export const ticketMutation = z.object({
  title: z.string().min(16).max(64),
  description: z.string().min(32).max(128),
  mode: z.boolean().optional(),
  quota: z.number().min(1).max(10_000_000),
  image: z
    .instanceof(File)
    .nullable()
    .or(z.literal(null))
    .refine((file) => {
      if (file === null) return true;
      return file.size <= MAX_FILE_SIZE;
    })
    .refine((file) => {
      if (file === null) return true;
      const ext = file.name.split(".").pop()?.toLowerCase();
      return ext && ACCEPTED_EXTENSIONS.includes(ext);
    }),
});

export type ticketMutationType = z.infer<typeof ticketMutation>;
