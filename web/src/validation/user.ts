import z from "zod";
import { ACCEPTED_EXTENSIONS, MAX_FILE_SIZE } from "./ticket";

export const userMutation = z.object({
  name: z.string().min(1).max(36),
  username: z.string().min(1).max(36),
  avatar: z
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

export type userMutationType = z.infer<typeof userMutation>;
