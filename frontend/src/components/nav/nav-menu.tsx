"use client";

import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerDescription,
  DrawerFooter,
  DrawerHeader,
  DrawerOverlay,
  DrawerPortal,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer";
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
import { NavSearch } from "./nav-search";
import { User } from "@/lib/schema/user";
import { usePathname, useRouter } from "next/navigation";
import { useMediaQuery } from "@/hooks/use-media-query";
import { MenuIcon } from "lucide-react";
import { ScrollArea } from "../ui/scroll-area";
import { NavMobile } from "./nav-mobile";

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
  const isDesktop = useMediaQuery("(min-width: 768px)");
  const router = useRouter();

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
              {isDesktop ? (
                <>
                  <NavigationMenuItem>
                    <NavSearch />
                  </NavigationMenuItem>
                  <NavigationMenuItem>
                    {user.data && <NoteCreateDialog />}
                  </NavigationMenuItem>
                </>
              ) : (
                <NavMobile isLogin={!user.error} />
              )}
              <NavigationMenuItem>
                <SetTheme />
              </NavigationMenuItem>
              {user.data ? (
                <NavigationMenuItem>
                  <UserMenu user={user.data} />
                </NavigationMenuItem>
              ) : (
                isDesktop && (
                  <NavigationMenuItem>
                    <Button
                      onClick={() => router.push("/login")}
                      className="w-24"
                    >
                      Login
                    </Button>
                  </NavigationMenuItem>
                )
              )}
            </NavigationMenuList>
          </NavigationMenu>
        </div>
      </div>
    );
  }
};

export { NavMenu };
