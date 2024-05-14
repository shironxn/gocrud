"use client";

import { Button } from "@/components/ui/button";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Ellipsis } from "lucide-react";
import { NavSearch } from "./nav-search";
import { NoteCreateDialog } from "../note-drawer";
import { useRouter } from "next/navigation";

export function NavMobile({ isLogin }: { isLogin: boolean }) {
  const router = useRouter();

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button variant="outline">
          <Ellipsis />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-80">
        <div className="grid gap-4">
          <NavSearch />
          {isLogin ? (
            <NoteCreateDialog />
          ) : (
            <Button onClick={() => router.push("/login")}>Login</Button>
          )}
        </div>
      </PopoverContent>
    </Popover>
  );
}
