"use client";

import { Input } from "./ui/input";
import Link from "next/link";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  NavigationMenu,
  NavigationMenuContent,
  NavigationMenuIndicator,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  NavigationMenuTrigger,
  NavigationMenuViewport,
  navigationMenuTriggerStyle,
} from "@/components/ui/navigation-menu";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button } from "./ui/button";
import {
  UserRound,
  Github,
  Smile,
  LogOut,
  SquareArrowOutUpRight,
  Search,
} from "lucide-react";
import React from "react";

const menu: { title: string; href: string; icon: JSX.Element }[] = [
  {
    title: "Profile",
    href: "/profile",
    icon: <UserRound />,
  },
  {
    title: "Github",
    href: "https://github.com/shironxn/gocrud",
    icon: <Github />,
  },
  {
    title: "Support",
    href: "https://saweria.co/shironxn",
    icon: <Smile />,
  },
  {
    title: "Api",
    href: "https://github.com/shironxn/gocrud",
    icon: <SquareArrowOutUpRight />,
  },
];

export default function Navbar() {
  return (
    <div className="shadow-md p-4 flex justify-between items-center">
      <div>
        <Link href={"/"}>
          <h1 className="scroll-m-20 text-2xl font-bold tracking-tight">
            NotesLand
          </h1>
        </Link>
      </div>
      <div>
        <NavigationMenu>
          <NavigationMenuList className="flex gap-2">
            <NavigationMenuItem>
              <div className="flex h-10 items-center rounded-md border border-input bg-white pl-3 text-sm ring-offset-background focus-within:ring-1 focus-within:ring-ring focus-within:ring-offset-2">
                <Search className="absolute pointer-events-none h-5" />
                <input
                  className="w-full p-2 pl-8 placeholder:text-muted-foreground focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-50"
                  placeholder="Search"
                />
              </div>
            </NavigationMenuItem>
            <NavigationMenuItem>
              <Button className="w-24">Create</Button>
            </NavigationMenuItem>
            <NavigationMenuItem>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Avatar className="cursor-pointe">
                    <AvatarImage src="https://github.com/shadcn.png" />
                    <AvatarFallback>CN</AvatarFallback>
                  </Avatar>
                </DropdownMenuTrigger>
                <DropdownMenuContent className="mr-5 mt-5 w-56">
                  <DropdownMenuItem disabled className="font-bold">
                    shironxn
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                  {menu.map((item, i) => (
                    <Link href={item.href} target="_blank" key={i}>
                      <DropdownMenuItem>
                        {item.icon}
                        {item.title}
                      </DropdownMenuItem>
                    </Link>
                  ))}
                  <DropdownMenuSeparator />
                  <DropdownMenuItem>
                    <LogOut className="h-4 w-4" />
                    Log out
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </NavigationMenuItem>
          </NavigationMenuList>
        </NavigationMenu>
      </div>
    </div>
  );
}
