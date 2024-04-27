"use client";

import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";

const NotePagination = ({
  currentPage,
  totalPage,
}: {
  currentPage: number;
  totalPage: number;
}) => {
  const searchParams = useSearchParams();
  const pathname = usePathname();
  const { replace } = useRouter();

  const [page, setPage] = useState(currentPage);

  useEffect(() => {
    const params = new URLSearchParams(searchParams);
    if (page !== 1) {
      params.set("page", page.toString());
    } else {
      params.delete("page");
    }
    replace(`${pathname}?${params.toString()}`);
  }, [page]);

  return (
    <Pagination>
      <PaginationContent>
        <PaginationItem>
          <PaginationPrevious
            onClick={() => setPage((prevPage) => prevPage - 1)}
            aria-disabled={page <= 1}
            tabIndex={page <= 1 ? -1 : undefined}
            className={
              page <= 1 ? "pointer-events-none opacity-50" : "cursor-pointer"
            }
          />
        </PaginationItem>
        <PaginationItem>
          <PaginationNext
            onClick={() => setPage((prevPage) => prevPage + 1)}
            aria-disabled={page >= totalPage}
            tabIndex={page >= totalPage ? totalPage : undefined}
            className={
              page >= totalPage
                ? "pointer-events-none opacity-50"
                : "cursor-pointer"
            }
          />
        </PaginationItem>
      </PaginationContent>
    </Pagination>
  );
};

export { NotePagination };
