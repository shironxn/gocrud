import Image from "next/image";
import { Navbar } from "@/components/navbar";
import { AspectRatio } from "@/components/ui/aspect-ratio";
import { NoteCard } from "@/components/note-card";

export default async function Page() {
  return (
    <div className="min-h-screen justify-center">
      <Navbar />
      <div className="w-full">
        <AspectRatio ratio={14 / 2}>
          <Image src="/header.jpg" alt="Image" className="object-cover" fill />
        </AspectRatio>
      </div>
      <div className="container mx-auto py-8">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <NoteCard />
        </div>
      </div>
    </div>
  );
}
