import "@/app/globals.css";
import { NoteSide } from "@/components/note/note-side";

export default function NoteLayout({
  children,
  params,
}: {
  children: React.ReactNode;
  params: { slug: string };
}) {
  return (
    <div className="flex min-h-screen gap-8">
      {/* <aside className="sticky top-0 h-[calc(100vh-theme(spacing.16))] w-40 overflow-y-auto bg-green-200"> */}
      <NoteSide id={params.slug} />
      {/* </aside> */}

      {children}
    </div>
  );
}
