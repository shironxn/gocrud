import { GetUserByName } from "@/actions/user";
import { GetNotes } from "@/actions/note";
import { NoteCard } from "@/components/note/note-card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { NotePagination } from "@/components/note/note-pagination";

export default async function Page({
  params,
  searchParams,
}: {
  params: { slug: string };
  searchParams?: {
    search?: string;
    page?: string;
  };
}) {
  const search = searchParams?.search || "";
  const page = Number(searchParams?.page) || 1;

  const user = await GetUserByName(params.slug);
  const notes = await GetNotes({
    search: search,
    user_id: user?.id,
    page: page,
    visibility: "public",
  });

  return (
    <div className="space-y-8">
      <div className="flex items-center justify-center mt-8">
        <div className="flex flex-col items-center">
          <Avatar className="h-24 w-24 md:h-32 md:w-32">
            <AvatarImage src={user.avatar_url} />
            <AvatarFallback>
              {user.name.slice(0, 2).toUpperCase()}
            </AvatarFallback>
          </Avatar>
          <div className="justify-center text-center items-center mt-8">
            <h3 className="scroll-m-20 text-2xl font-semibold tracking-tight">
              {user.name || "Unknown User"}
            </h3>
            <p className="leading-7 mt-4">{user.bio}</p>
          </div>
        </div>
      </div>
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
