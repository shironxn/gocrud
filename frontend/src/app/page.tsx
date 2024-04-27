"use client";

import Image from "next/image";
import { Navbar } from "@/components/navbar";
import { AspectRatio } from "@/components/ui/aspect-ratio";
import { LoadingCard, NoteCard } from "@/components/note-card";
import { NotePagination } from "@/components/note-pagination";
import useAxios from "axios-hooks";
import { toast } from "@/components/ui/use-toast";

export default function Page({
  searchParams,
}: {
  searchParams?: {
    search?: string;
    page?: string;
  };
}) {
  const search = searchParams?.search || "";
  const currentPage = Number(searchParams?.page) || 1;

  const [{ data, loading, error }] = useAxios({
    url: `/notes?visibility="public"&title=${search}&page=${currentPage}&limit=6&order=desc`,
    method: "GET",
    baseURL: process.env.NEXT_PUBLIC_API_URL,
    withCredentials: true,
  });

  if (error) {
    toast({
      title: "Uh oh! Something went wrong.",
      description: error.response?.data.message || "An unknown error occurred",
    });
  }

  return (
    <div className="min-h-screen justify-center">
      <Navbar />
      <div className="w-full">
        <AspectRatio ratio={14 / 2}>
          <Image src="/header.jpg" alt="Image" className="object-cover" fill />
        </AspectRatio>
      </div>
      {data?.data.notes ? (
        <div className="container mx-auto py-8 space-y-8">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {data && <NoteCard data={data.data.notes} />}
          </div>
          {Number(data.data?.metadata?.total_page) > 1 && (
            <NotePagination
              currentPage={currentPage}
              totalPage={Number(data.data?.metadata?.total_page)}
            />
          )}
        </div>
      ) : (
        <LoadingCard />
      )}
    </div>
  );
}
