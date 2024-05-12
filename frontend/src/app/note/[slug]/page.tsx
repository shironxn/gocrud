import { GetNotes, GetNotesByID } from "@/actions/note";

interface queryParams {
  page?: number;
}

export default async function Page({ params }: { params: { slug: string } }) {
  const note = await GetNotesByID(params.slug);
  return (
    <div className="w-3/4">
      <p className="leading-7 [&:not(:first-child)]:mt-6">{note.content}</p>
    </div>
  );
}
