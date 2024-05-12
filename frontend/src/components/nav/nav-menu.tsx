"use client";

import Link from "next/link";
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuList,
} from "@/components/ui/navigation-menu";
import { Button } from "@/components/ui/button";
import React from "react";
import { NoteCreateDialog } from "../note-drawer";
import { UserMenu } from "./user-menu";
import { SetTheme } from "./set-theme";
import { SearchNav } from "./search";
import { User } from "@/lib/schema/user";
import { usePathname } from "next/navigation";

const NavMenu = ({
  user,
}: {
  user:
    | {
        error: string;
        data?: User;
      }
    | {
        data: User;
        error?: string;
      };
}) => {
  const pathname = usePathname();

  if (!pathname.includes("/login") && !pathname.includes("/register")) {
    return (
      <div className="py-4 flex justify-between items-center">
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
                <SearchNav />
              </NavigationMenuItem>
              <NavigationMenuItem>
                {user.data && <NoteCreateDialog />}
              </NavigationMenuItem>
              <NavigationMenuItem>
                <SetTheme />
              </NavigationMenuItem>
              {user.data ? (
                <NavigationMenuItem>
                  <UserMenu user={user.data} />
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
  }
};

export { NavMenu };
