"use client";

import Image from "next/image";
import Navbar from "@/components/navbar";
import { AspectRatio } from "@/components/ui/aspect-ratio";
import { Avatar, AvatarImage } from "@/components/ui/avatar";
import { AvatarFallback } from "@radix-ui/react-avatar";
import { Button } from "@/components/ui/button";

export default function Home() {
  return (
    <div className="items-center justify-center">
      <Navbar />
      <div className="w-full mb-5">
        <AspectRatio ratio={14 / 2}>
          <Image src="/header.jpg" alt="Image" className="object-cover" fill />
        </AspectRatio>
      </div>
      <div className="mb-5 pl-5 ">
        <h3 className="scroll-m-20 text-2xl font-semibold tracking-tight">
          Public Notes
        </h3>
        <div className="flex overflow-x-auto">
          <div className="grid grid-cols-1 gap-4 grid-flow-col">
            <div className="flex gap-5 justify-between items-start px-6 py-4 font-medium bg-white rounded-lg border border-solid shadow-md border-zinc-200">
              <Avatar>
                <AvatarImage src="https://github.com/shadcn.png" />
                <AvatarFallback>CN</AvatarFallback>
              </Avatar>
              <div>
                <div className="text-sm font-bold leading-none">@shironxn</div>
                <div className="text-sm font-medium leading-non mt-2">
                  Jawa on top ilmu padi king 🔥
                </div>
                <p className="text-sm text-muted-foreground">5 minutes ago</p>
              </div>
              <Button>Read</Button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
