"use client";

import { AspectRatio } from "@/components/ui/aspect-ratio";
import Image from "next/image";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { toast } from "@/components/ui/use-toast";
import Link from "next/link";
import useAxios from "axios-hooks";

const FormSchema = z.object({
  email: z.string().email(),
  password: z.string(),
});

export default function Login() {
  const form = useForm<z.infer<typeof FormSchema>>({
    resolver: zodResolver(FormSchema),
  });

  const [{ data, loading, error }, executePost] = useAxios(
    {
      url: "http://127.0.0.1:3000/api/v1/auth/login",
      method: "POST",
      withCredentials: true,
    },
    { manual: true }
  );

  function onSubmit(data: z.infer<typeof FormSchema>) {
    executePost({ data: data });
    toast({
      title: "You submitted the following values:",
      description: (
        <pre className="mt-2 w-[340px] rounded-md bg-slate-950 p-4">
          <code className="text-white">{JSON.stringify(data, null, 2)}</code>
        </pre>
      ),
    });
  }

  return (
    <div className="flex justify-center items-center min-h-screen">
      <div className="max-w-6xl flex flex-col space-y-5 p-5 md:flex-row w-full md:space-x-10">
        <div className="w-full md:w-1/2 flex justify-center items-center ">
          <AspectRatio ratio={1 / 1}>
            <Image
              src="/login.png"
              alt="Image"
              className="rounded-md object-cover"
              fill
            />
          </AspectRatio>
        </div>
        <div className="w-full space-y-6 md:w-1/2">
          <div>
            <h1 className="text-center scroll-m-20 text-5xl font-extrabold tracking-tight md:text-5xl">
              Login
            </h1>
          </div>
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
              <FormField
                control={form.control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Email</FormLabel>
                    <FormControl>
                      <Input placeholder="example@gmail.com" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="password"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Password</FormLabel>
                    <FormControl>
                      <Input placeholder="your password" {...field} />
                    </FormControl>
                    <FormMessage />
                    <FormDescription className="text-right cursor-pointer">
                      Forgot password?
                    </FormDescription>
                  </FormItem>
                )}
              />
              <Button type="submit" className="w-full mt-3">
                Submit
              </Button>
            </form>
            <p className="leading-7 [&:not(:first-child)]:mt-6 text-center">
              Not have an account?{" "}
              <span className="font-semibold cursor-pointer">
                <Link href="/register">Register</Link>
              </span>
            </p>
          </Form>
        </div>
      </div>
    </div>
  );
}
