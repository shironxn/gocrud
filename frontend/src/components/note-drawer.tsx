"use client";

import * as React from "react";

import { cn } from "@/lib/utils";
import { useMediaQuery } from "@/hooks/use-media-query";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer";

import { User, UserRequest, userRequestSchema } from "@/lib/schema/user";
import { EllipsisVertical, PenBox } from "lucide-react";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

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
import { LoadingButton } from "@/components/loading-button";
import { useForm } from "react-hook-form";
import {
  Note,
  NoteCreate,
  noteCreateSchema,
  noteUpdateSchema,
  NoteUpdate,
} from "@/lib/schema/note";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "./ui/use-toast";
import { useEffect } from "react";
import useAxios from "axios-hooks";
import { useRouter } from "next/navigation";
import { buttonVariants } from "./ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "./ui/textarea";
import { Input } from "./ui/input";
import { Button } from "@/components/ui/button";

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
  const [open, setOpen] = React.useState(false);
  const isDesktop = useMediaQuery("(min-width: 768px)");

  if (isDesktop) {
    return (
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogTrigger asChild>
          <Button>Create</Button>
        </DialogTrigger>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Create note</DialogTitle>
            <DialogDescription>Create a new note</DialogDescription>
          </DialogHeader>
          <NoteForm />
        </DialogContent>
      </Dialog>
    );
  }

  return (
    <Drawer open={open} onOpenChange={setOpen}>
      <DrawerTrigger asChild>
        <Button>Create</Button>
      </DrawerTrigger>
      <DrawerContent>
        <DrawerHeader className="text-left">
          <DrawerTitle>Create note</DrawerTitle>
          <DrawerDescription>Create a new note</DrawerDescription>
        </DrawerHeader>
        <NoteForm className="px-4" />
        <DrawerFooter className="pt-2">
          <DrawerClose asChild>
            <Button variant="outline">Cancel</Button>
          </DrawerClose>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  );
};

const NoteUpdateDialog = ({ note }: { note: Note }) => {
  const [open, setOpen] = React.useState(false);
  const isDesktop = useMediaQuery("(min-width: 768px)");

  if (isDesktop) {
    return (
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogTrigger asChild>
          <button className="w-full text-left hover:bg-accent hover:text-accent-foreground px-4 py-2 rounded-md">
            Edit
          </button>
        </DialogTrigger>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Edit note</DialogTitle>
            <DialogDescription>Edit a note</DialogDescription>
          </DialogHeader>
          <NoteForm note={note} />
        </DialogContent>
      </Dialog>
    );
  }

  return (
    <Drawer open={open} onOpenChange={setOpen}>
      <DrawerTrigger asChild>
        <button className="w-full text-left hover:bg-accent hover:text-accent-foreground px-4 py-2 rounded-md">
          Edit
        </button>
      </DrawerTrigger>
      <DrawerContent>
        <DrawerHeader className="text-left">
          <DrawerTitle>Edit note</DrawerTitle>
          <DrawerDescription>Edit a note</DrawerDescription>
        </DrawerHeader>
        <NoteForm className="px-4" note={note} />
        <DrawerFooter className="pt-2">
          <DrawerClose asChild>
            <Button variant="outline">Cancel</Button>
          </DrawerClose>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  );
};

function NoteForm({ className, note }: { className?: string; note?: Note }) {
  const form = useForm<NoteCreate | NoteCreate>({
    resolver: zodResolver(!note ? noteCreateSchema : noteUpdateSchema),
  });

  const [
    { data: createData, loading: createLoading, error: createError },
    executeCreate,
  ] = useAxios(
    {
      url: "/notes",
      method: "POST",
      withCredentials: true,
      baseURL: process.env.NEXT_PUBLIC_API_URL,
    },
    { manual: true }
  );

  const [
    { data: updateData, loading: updateLoading, error: updateError },
    executeUpdate,
  ] = useAxios(
    {
      url: `/notes/${note?.id}`,
      method: "PUT",
      withCredentials: true,
      baseURL: process.env.NEXT_PUBLIC_API_URL,
    },
    { manual: true }
  );

  const onSubmit = (data: NoteCreate | NoteUpdate) => {
    if (note) {
      executeUpdate({ data: data });
    } else {
      executeCreate({ data: data });
    }
  };

  useEffect(() => {
    if (createError) {
      toast({
        title: "Uh oh! Something went wrong.",
        description:
          createError.response?.data.message || "An unknown error occurred",
      });
    }
  }, [createError]);

  useEffect(() => {
    if (updateError) {
      toast({
        title: "Uh oh! Something went wrong.",
        description:
          updateError.response?.data.message || "An unknown error occurred",
      });
    }
  }, [updateError]);

  useEffect(() => {
    if (createData) {
      window.location.reload();
    }
  }, [createData]);

  useEffect(() => {
    if (updateData) {
      window.location.reload();
    }
  }, [updateData]);

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className={cn("grid items-start gap-4", className)}
      >
        <div className="grid gap-4">
          {inputList.map((item, i) => (
            <div key={i}>
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
                          defaultValue={note?.visibility}
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
                        <Input
                          placeholder={item.placeholder}
                          defaultValue={note && (note[item.name] as string)}
                          {...field}
                        />
                      )}
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
          ))}
          <FormField
            control={form.control}
            name="content"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Content</FormLabel>
                <FormControl>
                  <Textarea
                    placeholder="Share your thoughts, ideas, or stories here..."
                    defaultValue={note?.content}
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>
        {!note ? (
          <LoadingButton type="submit" loading={createLoading}>
            Create
          </LoadingButton>
        ) : (
          <LoadingButton type="submit" loading={updateLoading}>
            Update
          </LoadingButton>
        )}
      </form>
    </Form>
  );
}

const NoteDeleteAlert = ({ note }: { note: Note }) => {
  const [
    { data: deleteData, loading: deleteLoading, error: deleteError },
    executeDelete,
  ] = useAxios(
    {
      url: `/notes/${note?.id}`,
      method: "DELETE",
      withCredentials: true,
      baseURL: process.env.NEXT_PUBLIC_API_URL,
    },
    { manual: true }
  );

  const handleClick = () => {
    executeDelete();
  };

  useEffect(() => {
    if (deleteError) {
      toast({
        title: "Uh oh! Something went wrong.",
        description:
          deleteError.response?.data.message || "An unknown error occurred",
      });
    }
  }, [deleteError]);

  useEffect(() => {
    if (deleteData) {
      window.location.reload();
    }
  }, [deleteData]);

  return (
    <AlertDialog>
      <AlertDialogTrigger asChild>
        <button className="w-full text-left hover:bg-accent hover:text-accent-foreground px-4 py-2 rounded-md">
          Delete
        </button>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Are you absolutely sure?</AlertDialogTitle>
          <AlertDialogDescription>
            This action cannot be undone. This will permanently delete your
            account and remove your data from our servers.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction onClick={handleClick}>Continue</AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
};

const NoteMenu = ({ note }: { note: Note }) => {
  const [open, setOpen] = React.useState(false);
  return (
    <DropdownMenu open={open} onOpenChange={setOpen}>
      <DropdownMenuTrigger>
        <Button variant={"outline"}>Menu</Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <DropdownMenuLabel>Note Menu</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <NoteUpdateDialog note={note} />
        <NoteDeleteAlert note={note} />
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export { NoteCreateDialog, NoteUpdateDialog, NoteMenu };
