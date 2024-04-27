"use client";

import Link from "next/link";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuList,
} from "@/components/ui/navigation-menu";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Button, buttonVariants } from "./ui/button";
import {
  UserRound,
  Github,
  Smile,
  LogOut,
  SquareArrowOutUpRight,
  Search,
  Moon,
  Sun,
} from "lucide-react";
import React, { useEffect } from "react";
import { useTheme } from "next-themes";
import useAxios from "axios-hooks";
import { toast } from "./ui/use-toast";
import { usePathname, useRouter } from "next/navigation";
import { NoteCreateDialog } from "./note-drawer";
import { useSearchParams } from "next/navigation";

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

const Navbar = () => {
  const { setTheme } = useTheme();
  const searchParams = useSearchParams();
  const pathname = usePathname();
  const { replace } = useRouter();

  const [{ data: userData, error: userError }, refetchUser] = useAxios(
    {
      url: "/users/me",
      method: "GET",
      baseURL: process.env.NEXT_PUBLIC_API_URL,
      withCredentials: true,
    },
    { useCache: false }
  );

  const [
    { data: logoutData, loading: logoutLoading, error: logoutError },
    executeLogout,
  ] = useAxios(
    {
      url: "/auth/logout",
      method: "POST",
      baseURL: process.env.NEXT_PUBLIC_API_URL,
      withCredentials: true,
    },
    { manual: true }
  );

  useEffect(() => {
    if (logoutError) {
      toast({
        title: "Uh oh! Something went wrong.",
        description:
          logoutError.response?.data.message || "An unknown error occurred",
      });
    }
  }, [logoutError]);

  useEffect(() => {
    if (logoutData) {
      window.location.reload();
    }
  }, [logoutData]);

  const onClick = () => {
    executeLogout();
  };

  const handleSearch = (term: string) => {
    const params = new URLSearchParams(searchParams);
    params.delete("page");
    if (term) {
      params.set("search", term);
    } else {
      params.delete("search");
    }
    replace(`${pathname}?${params.toString()}`);
  };

  return (
    <div className="shadow-md p-4 flex justify-between items-center">
      <div>
        <Link href={"/"}>
          <h1 className="scroll-m-20 text-2xl font-bold tracking-tight">
            BlankNotes
          </h1>
        </Link>
      </div>
      <div>
        <NavigationMenu>
          <NavigationMenuList className="flex gap-2">
            <NavigationMenuItem>
              <div className="flex h-10 items-center rounded-md border border-input pl-3 text-sm ring-offset-background focus-within:ring-1 focus-within:ring-ring focus-within:ring-offset-2">
                <Search className="absolute pointer-events-none h-5" />
                <input
                  className="w-full bg-transparent p-2 pl-8 placeholder:text-muted-foreground focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-50"
                  placeholder="Search"
                  onChange={(e) => {
                    handleSearch(e.target.value);
                  }}
                  defaultValue={searchParams.get("query")?.toString()}
                />
              </div>
            </NavigationMenuItem>
            <NavigationMenuItem>
              {userData && !userError && <NoteCreateDialog />}
            </NavigationMenuItem>
            <NavigationMenuItem>
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="outline" size="icon">
                    <Sun className="h-[1.2rem] w-[1.2rem] rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
                    <Moon className="absolute h-[1.2rem] w-[1.2rem] rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
                    <span className="sr-only">Toggle theme</span>
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end">
                  <DropdownMenuItem onClick={() => setTheme("light")}>
                    Light
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => setTheme("dark")}>
                    Dark
                  </DropdownMenuItem>
                  <DropdownMenuItem onClick={() => setTheme("system")}>
                    System
                  </DropdownMenuItem>
                </DropdownMenuContent>
              </DropdownMenu>
            </NavigationMenuItem>
            {userData && !userError ? (
              <NavigationMenuItem>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Avatar className="h-12 w-12 cursor-pointer">
                      <AvatarImage src={userData?.data.avatar_url} />
                      <AvatarFallback>
                        {userData?.data.name.slice(0, 2).toUpperCase()}
                      </AvatarFallback>
                    </Avatar>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent className="mr-5 mt-5 w-56">
                    <DropdownMenuItem disabled className="font-bold">
                      {userData?.data?.name}
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    {menu.map((item, i) => (
                      <Link href={item.href} key={i}>
                        <DropdownMenuItem className="gap-x-3">
                          <div>{item.icon}</div>
                          <div>{item.title}</div>
                        </DropdownMenuItem>
                      </Link>
                    ))}
                    <DropdownMenuSeparator />
                    <DropdownMenuItem
                      onClick={onClick}
                      disabled={logoutLoading}
                      className="gap-x-3"
                    >
                      <LogOut />
                      <div>Log out</div>
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </NavigationMenuItem>
            ) : (
              <Link href={"/login"}>
                <Button className="w-24">Login</Button>
              </Link>
            )}
          </NavigationMenuList>
        </NavigationMenu>
      </div>
    </div>
  );
};

export { Navbar };
