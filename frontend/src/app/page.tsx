"use client";

import Image from "next/image";
import { Navbar } from "@/components/navbar";
import { AspectRatio } from "@/components/ui/aspect-ratio";
import { LoadingCard, NoteCard } from "@/components/note-card";
import { NotePagination } from "@/components/pagination";
import useAxios from "axios-hooks";
import { toast } from "@/components/ui/use-toast";

export default function Page({
  searchParams,
}: {
  searchParams?: {
    query?: string;
    page?: string;
  };
}) {
  const query = searchParams?.query || "";
  const currentPage = Number(searchParams?.page) || 1;

  const [{ data, loading, error }] = useAxios({
    url: `/notes?title=${query}&page=${currentPage}&limit=6&order=desc`,
    method: "GET",
    baseURL: process.env.NEXT_PUBLIC_API_URL,
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
      {data?.data.notes && (
        <div className="container mx-auto py-8 space-y-8">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {!loading ? <NoteCard data={data} /> : <LoadingCard />}
          </div>
          {Number(data.data?.metadata?.total_page) > 1 && (
            <NotePagination
              currentPage={currentPage}
              totalPage={Number(data.data?.metadata?.total_page)}
            />
          )}
        </div>
      )}
    </div>
  );
}
