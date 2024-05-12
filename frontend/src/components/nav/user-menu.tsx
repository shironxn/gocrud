"use client";

import Link from "next/link";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  UserRound,
  Github,
  Smile,
  LogOut,
  SquareArrowOutUpRight,
} from "lucide-react";
import React from "react";
import { toast } from "../ui/use-toast";
import { useRouter } from "next/navigation";
import { User } from "@/lib/schema/user";
import { Logout } from "@/actions/auth";

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

const UserMenu = ({ user }: { user: User }) => {
  const router = useRouter();

  const onClick = async () => {
    const error = await Logout();
    if (error) {
      toast({
        title: "Uh oh! Something went wrong.",
        description: error,
      });
    } else {
      router.refresh();
    }
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Avatar className="h-12 w-12 cursor-pointer">
          <AvatarImage src={user.avatar_url} />
          <AvatarFallback>
            {user.name?.slice(0, 2).toUpperCase()}
          </AvatarFallback>
        </Avatar>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="mr-5 mt-5 w-56">
        <DropdownMenuItem disabled className="font-bold">
          {user?.name}
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
        <DropdownMenuItem onClick={onClick} className="gap-x-3">
          <LogOut />
          <div>Log out</div>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export { UserMenu };
