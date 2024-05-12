"use client";

import { Search } from "lucide-react";
import React from "react";
import { usePathname, useRouter } from "next/navigation";
import { useSearchParams } from "next/navigation";

const SearchNav = () => {
  const searchParams = useSearchParams();
  const pathname = usePathname();
  const router = useRouter();

  const handleSearch = (term: string) => {
    const params = new URLSearchParams(searchParams);
    params.delete("page");
    if (term) {
      params.set("search", term);
    } else {
      params.delete("search");
    }
    router.replace(`${pathname}?${params.toString()}`);
  };

  return (
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
  );
};

export { SearchNav };
