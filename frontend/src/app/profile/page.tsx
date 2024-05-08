"use client";

import { Navbar } from "@/components/navbar";
import { LoadingCard, NoteCard } from "@/components/note-card";
import { NotePagination } from "@/components/note-pagination";
import useAxios from "axios-hooks";
import { toast } from "@/components/ui/use-toast";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { UserUpdateDrawer } from "@/components/user-drawer";

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

  const [{ data: userData }] = useAxios(
    {
      url: "/users/me",
      method: "GET",
      baseURL: process.env.NEXT_PUBLIC_API_URL,
      withCredentials: true,
    },
    { useCache: false }
  );

  const [{ data, loading, error }] = useAxios({
    url: `/notes?user_id=${userData?.data.id || 0}&title=${search}&author=${
      userData?.data.name
    }&page=${currentPage}&limit=6&order=desc`,
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
    <div>
      <Navbar />
      {data && (
        <div className="flex items-center justify-center mt-8">
          <div className="flex flex-col items-center">
            <Avatar className="h-24 w-24 md:h-32 md:w-32">
              <AvatarImage src={userData?.data.avatar_url} />
              <AvatarFallback>
                {userData?.data.name.slice(0, 2).toUpperCase()}
              </AvatarFallback>
            </Avatar>
            <div className="justify-center text-center items-center mt-8">
              <h3 className="scroll-m-20 text-2xl font-semibold tracking-tight">
                {userData?.data?.name || "Unknown User"}
              </h3>
              <p className="leading-7 mt-4">{userData?.data?.bio}</p>
              <button className="relative scale-75">
                {<UserUpdateDrawer user={userData?.data} />}
              </button>
            </div>
          </div>
        </div>
      )}

      {data?.data.notes ? (
        <div className="container mx-auto py-8 space-y-8">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {data && <NoteCard data={data.data.notes} />}
          </div>
          {Number(data.data?.metadata?.total_pages) > 1 && (
            <NotePagination
              currentPage={currentPage}
              totalPages={Number(data.data?.metadata?.total_pages)}
            />
          )}
        </div>
      ) : (
        <LoadingCard />
      )}
    </div>
  );
}
