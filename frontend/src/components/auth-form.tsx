"use client";

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
import Link from "next/link";
import { LoadingButton } from "@/components/loading-button";
import { PasswordInput } from "@/components/password-input";
import { useForm } from "react-hook-form";
import {
  AuthLogin,
  authLoginSchema,
  AuthRegister,
  authRegisterSchema,
} from "@/lib/schema/auth";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "./ui/use-toast";
import { useEffect } from "react";
import useAxios from "axios-hooks";
import { useRouter } from "next/navigation";

const LoginForm = () => {
  const router = useRouter();

  const form = useForm<AuthLogin>({
    resolver: zodResolver(authLoginSchema),
  });

  const [{ data, loading, error }, executeLogin] = useAxios(
    {
      url: "/auth/login",
      method: "POST",
      baseURL: process.env.NEXT_PUBLIC_API_URL,
      withCredentials: true,
    },
    { manual: true }
  );

  const onSubmit = async (data: AuthLogin) => {
    executeLogin({ data: data });
  };

  useEffect(() => {
    if (error) {
      toast({
        title: "Uh oh! Something went wrong.",
        description:
          error.response?.data.message || "An unknown error occurred",
      });
    }
  }, [error]);

  useEffect(() => {
    if (data) {
      toast({
        title: "Success",
        description: data?.message,
      });
      router.push("/");
    }
  }, [data]);

  return (
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
                <PasswordInput placeholder="your password" {...field} />
              </FormControl>
              <FormMessage />
              <FormDescription className="text-right cursor-pointer">
                Forgot password?
              </FormDescription>
            </FormItem>
          )}
        />
        <LoadingButton loading={loading} type="submit" className="w-full mt-3">
          Submit
        </LoadingButton>
      </form>
      <p className="leading-7 [&:not(:first-child)]:mt-6 text-center">
        Not have an account?{" "}
        <span className="font-semibold cursor-pointer">
          <Link href="/register">Register</Link>
        </span>
      </p>
    </Form>
  );
};

const RegisterForm = () => {
  const router = useRouter();

  const form = useForm<AuthRegister>({
    resolver: zodResolver(authRegisterSchema),
  });

  const [{ data, loading, error }, executePost] = useAxios(
    {
      url: "/auth/register",
      method: "POST",
      baseURL: process.env.NEXT_PUBLIC_API_URL,
    },
    { manual: true }
  );

  function onSubmit(dataSubmit: AuthRegister) {
    executePost({ data: dataSubmit });
  }

  useEffect(() => {
    if (error) {
      toast({
        title: "Uh oh! Something went wrong.",
        description:
          error.response?.data.message || "An unknown error occurred",
      });
    }
  }, [error]);

  useEffect(() => {
    if (data) {
      toast({
        title: "Success",
        description: data?.message,
      });
      router.push("/login");
    }
  }, [data]);

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <FormControl>
                <Input placeholder="agus" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
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
                <PasswordInput placeholder="your password" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <LoadingButton loading={loading} type="submit" className="w-full mt-3">
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

export { LoginForm, RegisterForm };
