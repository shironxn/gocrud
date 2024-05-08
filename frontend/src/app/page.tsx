import Image from "next/image";
import { Navbar } from "@/components/navbar";
import { AspectRatio } from "@/components/ui/aspect-ratio";
import { LoadingCard, NoteCard } from "@/components/note-card";
import { NotePagination } from "@/components/note-pagination";
import useAxios from "axios-hooks";
import { toast } from "@/components/ui/use-toast";
import { GetNotes } from "@/actions/note";
import { Header } from "@/components/header";

interface queryParams {
  search?: string;
  page?: number;
}

export default async function Page({
  searchParams,
}: {
  searchParams?: queryParams;
}) {
  const search = searchParams?.search || "";
  const currentPage = Number(searchParams?.page) || 1;
  const data = await GetNotes({ page: currentPage, search: search });

  return (
    <div className="min-h-screen justify-center">
      <Navbar />
      <Header />
      {data?.notes ? (
        <div className="container mx-auto py-8 space-y-8">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {data && <NoteCard data={data.notes} />}
          </div>
          {Number(data?.metadata?.total_pages) > 1 && (
            <NotePagination
              currentPage={currentPage}
              totalPages={Number(data?.metadata?.total_pages)}
            />
          )}
        </div>
      ) : (
        <LoadingCard />
      )}
    </div>
  );
}
