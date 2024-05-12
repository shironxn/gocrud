import { GetNotes, GetNotesByID } from "@/actions/note";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { AspectRatio } from "@/components/ui/aspect-ratio";
import Image from "next/image";
import { Separator } from "../ui/separator";
import { NoteCard } from "./note-card";

const NoteSide = async ({ id }: { id: string }) => {
  const note = await GetNotesByID(id);
  const notes = await GetNotes({ author: note.author.name });
  return (
    <aside className="w-1/4 h-screen space-y-3 sticky top-0 overflow-auto no-scrollbar">
      <div>
        <h2 className="scroll-m-20 text-3xl font-semibold tracking-tight first:mt-0">
          {note.title}
        </h2>
        <p className="leading-7 pb-3">{note.description}</p>
        <AspectRatio ratio={16 / 9}>
          <Image
            src={note.cover_url}
            alt="Image"
            className="rounded-md object-cover"
            fill
          />
        </AspectRatio>
      </div>

      <div>
        <div className="flex text-center items-center space-x-4">
          <Avatar className="h-10 w-10">
            <AvatarImage src={note.author.avatar_url} />
            <AvatarFallback>
              {note.author.name.slice(0, 2).toUpperCase()}
            </AvatarFallback>
          </Avatar>
          <small className="text-sm font-medium leading-none">
            @{note.author.name}
          </small>
        </div>
        <div className="content-end">
          <p className="text-sm text-muted-foreground">
            <small className="text-sm text-muted-foreground">
              Created at{" "}
              {new Date(note.created_at).toLocaleDateString("en-DB", {
                day: "numeric",
                month: "long",
                year: "numeric",
              })}
            </small>
          </p>
          <p className="text-sm text-muted-foreground">
            Last update:{" "}
            <small className="text-sm text-muted-foreground">
              {new Date(note.updated_at).toLocaleDateString("en-DB", {
                day: "numeric",
                month: "long",
                year: "numeric",
              })}
            </small>
          </p>
        </div>
      </div>
      <div className="w-full">
        <Separator />
      </div>
      <div>
        <h3 className="scroll-m-20 text-2xl font-semibold tracking-tight">
          Read More
        </h3>
      </div>
      <div className="space-y-8">
        <NoteCard data={notes} />
      </div>
    </aside>
  );
};

export { NoteSide };
