import { NoteCard } from "@/components/note/note-card";
import { GetNotes } from "@/actions/note";
import { Header } from "@/components/header";
import { NotePagination } from "@/components/note/note-pagination";
import { number } from "zod";

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
  const page = Number(searchParams?.page) || 1;
  const notes = await GetNotes({
    page: page,
    search: search,
    visibility: "public",
  });

  return (
    <div className="min-h-screen justify-center space-y-8">
      <Header />
      <div className="flex-grow mx-auto">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8">
          <NoteCard data={notes} />
        </div>
        {Number(notes?.metadata?.total_pages) > 1 && (
          <div className="mt-8">
            <NotePagination
              currentPage={notes.metadata.page}
              totalPages={Number(notes?.metadata?.total_pages)}
            />
          </div>
        )}
      </div>
    </div>
  );
}
