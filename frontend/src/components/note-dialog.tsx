"use client";

import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";

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
import { AuthLogin, authLoginSchema } from "@/lib/schema/auth";
import { Note, NoteRequest, noteRequestSchema } from "@/lib/schema/note";
import { Success } from "@/lib/schema/response";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "./ui/use-toast";
import { useEffect } from "react";
import useAxios from "axios-hooks";
import { useRouter } from "next/navigation";
import { Button, buttonVariants } from "./ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "./ui/textarea";

type NameOptions =
  | "title"
  | "description"
  | "cover_url"
  | "content"
  | "visibility";
const inputList: { name: NameOptions; label: string; placeholder: string }[] = [
  {
    name: "title",
    label: "Title",
    placeholder: "Give your note a catchy title",
  },
  {
    name: "description",
    label: "Description",
    placeholder: "Briefly describe your note",
  },
  {
    name: "cover_url",
    label: "Cover URL",
    placeholder: "Paste a link to a cool cover image",
  },
  {
    name: "visibility",
    label: "Visibility",
    placeholder: "Choose who can see your note",
  },
];

const NoteCreateDialog = () => {
  const router = useRouter();

  const form = useForm<NoteRequest>({
    resolver: zodResolver(noteRequestSchema),
  });

  const [{ data, loading, error }, executeCreate] = useAxios(
    {
      url: "/notes",
      method: "POST",
      baseURL: process.env.NEXT_PUBLIC_API_URL,
      withCredentials: true,
    },
    { manual: true }
  );

  const onSubmit = async (data: NoteRequest) => {
    executeCreate({ data: data });
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
      const fields: string[] = Object.keys(form.getValues());
      fields.forEach((field: string) => {
        form.setValue(field as NameOptions, "");
      });
    }
  }, [data]);

  return (
    <AlertDialog>
      <AlertDialogTrigger className={buttonVariants()}>
        Create
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Create note</AlertDialogTitle>
          <AlertDialogDescription>Create a new note</AlertDialogDescription>
        </AlertDialogHeader>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            <div className="flex flex-wrap md:-mx-2">
              {inputList.map((item, i) => (
                <div key={i} className="w-full md:w-1/2 md:px-2 mb-6">
                  <FormField
                    control={form.control}
                    name={item.name}
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{item.label}</FormLabel>
                        <FormControl>
                          {item.name === "visibility" ? (
                            <Select
                              onValueChange={field.onChange}
                              defaultValue={field.value}
                            >
                              <SelectTrigger>
                                <SelectValue placeholder={item.placeholder} />
                              </SelectTrigger>
                              <SelectContent>
                                <SelectItem value="public">public</SelectItem>
                                <SelectItem value="private">private</SelectItem>
                              </SelectContent>
                            </Select>
                          ) : (
                            <Input placeholder={item.placeholder} {...field} />
                          )}
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
              ))}
              <div className="w-full px-2">
                <FormField
                  control={form.control}
                  name="content"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Content</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder="Share your thoughts, ideas, or stories here..."
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
            </div>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <AlertDialogAction type="submit">Create</AlertDialogAction>
            </AlertDialogFooter>
          </form>
        </Form>

        {/* <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
            <div className="flex flex-wrap">
              {inputList.map((item, i) => (
                <div key={i} className="w-1/2 p-2">
                  <FormField
                    control={form.control}
                    name={item.name}
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>{item.label}</FormLabel>
                        <FormControl>
                          <Input placeholder={item.placeholder} {...field} />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>
              ))}
              <div className="w-1/2 p-2">
                <FormField
                  control={form.control}
                  name="visibility"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Visibility</FormLabel>
                      <Select
                        onValueChange={field.onChange}
                        defaultValue={field.value}
                      >
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Choose who can see your note" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value="public">public</SelectItem>
                          <SelectItem value="private">private</SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
              <div className="w-full">
                <FormField
                  control={form.control}
                  name="content"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Content</FormLabel>
                      <FormControl>
                        <Textarea
                          placeholder="Share your thoughts, ideas, or stories here..."
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
            </div>
            <AlertDialogFooter>
              <AlertDialogCancel>Cancel</AlertDialogCancel>
              <LoadingButton loading={loading} type="submit">
                Create
              </LoadingButton>
            </AlertDialogFooter>
          </form>
        </Form> */}
      </AlertDialogContent>
    </AlertDialog>
  );
};

export { NoteCreateDialog };
