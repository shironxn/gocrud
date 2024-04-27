"use client";

import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { AspectRatio } from "./ui/aspect-ratio";
import Image from "next/image";
import { Note } from "@/lib/schema/note";
import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar";
import { Button } from "./ui/button";

import { Skeleton } from "./ui/skeleton";
import { NoteMenu, NoteUpdateDialog } from "./note-drawer";
import { usePathname } from "next/navigation";
import { useState } from "react";

export const LoadingCard = () => {
  return (
    <div className="container mx-auto py-8 space-y-8">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {Array.from({ length: 6 }).map((_, i: number) => (
          <Card key={i}>
            <CardHeader>
              <div className="w-full">
                <AspectRatio ratio={5 / 1}>
                  <Skeleton className="w-full h-full rounded-t-md" />
                </AspectRatio>
              </div>
            </CardHeader>
            <CardContent className="w-[248px] h-[104px]"></CardContent>
            <div className="px-6 pb-3">
              <Separator />
            </div>
            <CardFooter className="justify-between">
              <div className="flex text-center items-center space-x-4">
                <Avatar>
                  <Skeleton className="w-full h-full" />
                </Avatar>
                <small className="text-sm font-medium leading-none">
                  <Skeleton className="w-full h-full" />
                </small>
              </div>
              <Skeleton className="w-16 h-10" />
            </CardFooter>
          </Card>
        ))}
      </div>
    </div>
  );
};

export function NoteCard({ data }: { data: Note[] }) {
  const pathname = usePathname();
  const [imgError, setImgError] = useState(false);

  return (
    <>
      {data &&
        data.map((item: Note, i: number) => (
          <Card key={i}>
            <CardHeader>
              <div className="w-full">
                <AspectRatio ratio={5 / 1}>
                  <Image
                    src={imgError ? "/cover.jpg" : item.cover_url}
                    onError={() => setImgError(true)}
                    alt=""
                    className="object-cover rounded-t-md"
                    fill
                  />
                </AspectRatio>
              </div>
            </CardHeader>
            <CardContent className="md:h-[104px]">
              <div className="flex flex-col md:flex-row justify-between">
                <CardTitle className="scroll-m-20 text-xl font-semibold tracking-tigh line-clamp-2 md:w-2/4 lg:w-1/2">
                  {item.title}
                </CardTitle>
                <CardDescription>
                  {new Date(item.created_at).toLocaleDateString("en-DB", {
                    day: "numeric",
                    month: "long",
                    year: "numeric",
                  })}
                </CardDescription>
              </div>
              <p className="leading-7 line-clamp-2">{item.description}</p>
            </CardContent>
            <div className="px-6 pb-3">
              <Separator />
            </div>
            <CardFooter className="justify-between">
              <div className="flex text-center items-center space-x-4">
                <Avatar className="h-12 w-12">
                  <AvatarImage src={item.author.avatar_url} />
                  <AvatarFallback>
                    {item.author.name.slice(0, 2).toUpperCase()}
                  </AvatarFallback>
                </Avatar>
                <small className="text-sm font-medium leading-none">
                  @{item.author.name}
                </small>
              </div>
              <div>
                <Dialog>
                  <div className="space-x-2">
                    {pathname == "/profile" && <NoteMenu note={item} />}
                    <DialogTrigger asChild>
                      <Button variant="outline">Read</Button>
                    </DialogTrigger>
                  </div>
                  <DialogContent className="sm:max-w-md">
                    <DialogHeader>
                      <DialogTitle>{item.title}</DialogTitle>
                      <DialogDescription>
                        <p className="leading-7 [&:not(:first-child)]:mt-6">
                          {item.description}
                        </p>
                      </DialogDescription>
                      <div className="w-full">
                        <AspectRatio ratio={4 / 2}>
                          <Image
                            src={imgError ? "/cover.jpg" : item.cover_url}
                            onError={() => setImgError(true)}
                            alt=""
                            className="object-cover"
                            fill
                          />
                        </AspectRatio>
                      </div>
                    </DialogHeader>
                    <div className="flex items-center space-x-2">
                      <div className="grid flex-1 gap-2">
                        <p className="leading-7 [&:not(:first-child)]:mt-6">
                          {item.content}
                        </p>
                      </div>
                    </div>
                    <div className="w-full">
                      <Separator />
                    </div>
                    <DialogFooter>
                      <div className="justify-between flex items-center w-full">
                        <div className="flex text-center items-center space-x-4">
                          <Avatar className="h-12 w-12">
                            <AvatarImage src={item.author.avatar_url} />
                            <AvatarFallback>
                              {item.author.name.slice(0, 2).toUpperCase()}
                            </AvatarFallback>
                          </Avatar>
                          <small className="text-sm font-medium leading-none">
                            @{item.author.name}
                          </small>
                        </div>
                        <DialogDescription>
                          <p className="text-right">
                            {new Date(item.created_at).toLocaleDateString(
                              "en-DB",
                              {
                                day: "numeric",
                                month: "long",
                                year: "numeric",
                                hour: "numeric",
                                minute: "numeric",
                              }
                            )}
                          </p>
                        </DialogDescription>
                      </div>
                    </DialogFooter>
                  </DialogContent>
                </Dialog>
              </div>
            </CardFooter>
          </Card>
        ))}
    </>
  );
}
