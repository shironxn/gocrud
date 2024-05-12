"use client";

import * as React from "react";

import { cn } from "@/lib/utils";
import { useMediaQuery } from "@/hooks/use-media-query";
import { Button } from "@/components/ui/button";
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
import { Input } from "@/components/ui/input";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { User, UserRequest, userRequestSchema } from "@/lib/schema/user";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { PenBox } from "lucide-react";
import useAxios from "axios-hooks";
import { toast } from "./ui/use-toast";
import { useEffect } from "react";
import { LoadingButton } from "./loading-button";

type NameOptions = "name" | "email" | "password" | "bio" | "avatar_url";
const inputList: { name: NameOptions; label: string; placeholder: string }[] = [
  {
    name: "name",
    label: "Name",
    placeholder: "Give your note a catchy name",
  },
  {
    name: "email",
    label: "Email",
    placeholder: "Briefly describe your note",
  },
  {
    name: "password",
    label: "Password",
    placeholder: "Paste a link to a cool cover image",
  },
  {
    name: "bio",
    label: "Bio",
    placeholder: "Choose who can see your note",
  },
  {
    name: "avatar_url",
    label: "Avatar URL",
    placeholder: "Choose who can see your note",
  },
];

const UserUpdateDrawer = ({ user }: { user: User }) => {
  const [open, setOpen] = React.useState(false);
  const isDesktop = useMediaQuery("(min-width: 768px)");

  if (isDesktop) {
    return (
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogTrigger asChild>
          <PenBox />
        </DialogTrigger>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Edit profile</DialogTitle>
            <DialogDescription>
              Make changes to your profile here. Click save when you're done.
            </DialogDescription>
          </DialogHeader>
          <ProfileForm user={user} />
        </DialogContent>
      </Dialog>
    );
  }

  return (
    <Drawer open={open} onOpenChange={setOpen}>
      <DrawerTrigger asChild>
        <PenBox />
      </DrawerTrigger>
      <DrawerContent>
        <DrawerHeader className="text-left">
          <DrawerTitle>Edit profile</DrawerTitle>
          <DrawerDescription>
            Make changes to your profile here. Click save when you're done.
          </DrawerDescription>
        </DrawerHeader>
        <ProfileForm className="px-4" user={user} />
        <DrawerFooter className="pt-2">
          <DrawerClose asChild>
            <Button variant="outline">Cancel</Button>
          </DrawerClose>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  );
};

function ProfileForm({ className, user }: { className?: string; user: User }) {
  const form = useForm<UserRequest>({
    resolver: zodResolver(userRequestSchema),
    defaultValues: { bio: user.bio, avatar_url: user.avatar_url },
  });

  const [{ data, loading, error }, execute] = useAxios(
    {
      url: `/users/${user.id}`,
      method: "PUT",
      withCredentials: true,
      baseURL: process.env.NEXT_PUBLIC_API_URL,
    },
    { manual: true }
  );

  const onSubmit = (data: UserRequest) => {
    execute({ data: data });
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
      window.location.reload();
    }
  }, [data]);

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className={cn("grid items-start gap-4", className)}
      >
        <div className="grid gap-2">
          {inputList.map((item, i) => (
            <div key={i} className="w-full mb-6">
              <FormField
                control={form.control}
                name={item.name}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>{item.label}</FormLabel>
                    <FormControl>
                      <Input
                        placeholder={item.placeholder}
                        defaultValue={user[item.name] as string}
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
          ))}
        </div>
        <LoadingButton type="submit" loading={loading}>
          Save changes
        </LoadingButton>
      </form>
    </Form>
  );
}

export { UserUpdateDrawer };
