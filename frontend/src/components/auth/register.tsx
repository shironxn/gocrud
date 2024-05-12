"use client";

import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import Link from "next/link";
import { LoadingButton } from "@/components/loading-button";
import { PasswordInput } from "@/components/password-input";
import { useForm } from "react-hook-form";
import { AuthRegister, authRegisterSchema } from "@/lib/schema/auth";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "../ui/use-toast";
import { useRouter } from "next/navigation";
import { Register } from "@/actions/auth";

const RegisterForm = () => {
  const router = useRouter();

  const form = useForm<AuthRegister>({
    resolver: zodResolver(authRegisterSchema),
  });

  async function onSubmit(dataSubmit: AuthRegister) {
    const error = await Register(dataSubmit);
    if (error) {
      toast({
        title: "Uh oh! Something went wrong.",
        description: error,
      });
    } else {
      router.push("/login");
    }
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <div className="flex gap-x-2">
                <FormLabel>Name</FormLabel>
                <FormMessage className="text-xs" />
              </div>
              <FormControl>
                <Input required {...field} />
              </FormControl>
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="email"
          render={({ field }) => (
            <FormItem>
              <div className="flex gap-x-2">
                <FormLabel>Email</FormLabel>
                <FormMessage className="text-xs" />
              </div>
              <FormControl>
                <Input required type="email" {...field} />
              </FormControl>
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="password"
          render={({ field }) => (
            <FormItem>
              <div className="flex gap-x-2">
                <FormLabel>Password</FormLabel>
                <FormMessage className="text-xs" />
              </div>
              <FormControl>
                <PasswordInput required {...field} />
              </FormControl>
            </FormItem>
          )}
        />
        <LoadingButton
          loading={form.formState.isSubmitting}
          type="submit"
          className="w-full mt-3"
        >
          Submit
        </LoadingButton>
      </form>
      <p className="leading-7 [&:not(:first-child)]:mt-6 text-center">
        Already have an account?{" "}
        <span className="font-semibold cursor-pointer">
          <Link href="/login">Login</Link>
        </span>
      </p>
    </Form>
  );
};

export { RegisterForm };
