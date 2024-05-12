import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { UserUpdateDrawer } from "@/components/user-drawer";
import { GetUserMe } from "@/actions/user";
import { GetNotes } from "@/actions/note";
import { NoteCard } from "@/components/note/note-card";
import { toast } from "@/components/ui/use-toast";
import { NotePagination } from "@/components/note/note-pagination";

export default async function Page({
  searchParams,
}: {
  searchParams?: {
    search?: string;
    page?: string;
  };
}) {
  const search = searchParams?.search || "";
  const page = Number(searchParams?.page) || 1;

  const user = await GetUserMe();
  user.error &&
    toast({
      title: "Uh oh! Something went wrong.",
      description: user.error,
    });
  const notes = await GetNotes({
    search: search,
    user_id: user.data.id,
    page: page,
  });

  return (
    <div className="space-y-8">
      {user && (
        <div className="flex items-center justify-center mt-8">
          <div className="flex flex-col items-center">
            <Avatar className="h-24 w-24 md:h-32 md:w-32">
              <AvatarImage src={user.data.avatar_url} />
              <AvatarFallback>
                {user.data.name?.slice(0, 2).toUpperCase()}
              </AvatarFallback>
            </Avatar>
            <div className="justify-center text-center items-center mt-8">
              <h3 className="scroll-m-20 text-2xl font-semibold tracking-tight">
                {user.data.name || "Unknown User"}
              </h3>
              <p className="leading-7 mt-4">{user.data.bio}</p>
              <button className="relative scale-75">
                {<UserUpdateDrawer user={user.data} />}
              </button>
            </div>
          </div>
        </div>
      )}
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
