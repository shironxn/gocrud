"use server";

import { toast } from "@/components/ui/use-toast";
import { Note, NoteCreate, NoteUpdate } from "@/lib/schema/note";

interface queryParams {
  search?: string;
  page?: number;
}

const CreateNotes = async (note: NoteCreate) => {
  const res = await fetch(process.env.NEXT_PUBLIC_API_URL + `/notes`, {
    method: "POST",
    body: JSON.stringify(note),
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
  });
};

const GetNotes = async (query: queryParams) => {
  const res = await fetch(
    process.env.NEXT_PUBLIC_API_URL +
      `/notes?visibility="public"&title=${query.search}&page=${query.page}&limit=6&order=desc`
  );
  return res.json();
};

const UpdateNotes = async (note: NoteUpdate, id: string) => {
  const res = await fetch(process.env.NEXT_PUBLIC_API_URL + `/notes/${id}`, {
    method: "PUT",
    body: JSON.stringify(note),
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
  });
  return res.json();
};

const DeleteNotes = async (id: string) => {
  const res = await fetch(process.env.NEXT_PUBLIC_API_URL + `/notes/${id}`, {
    method: "DELETE",
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
    },
  });
  return res.json();
};

export { CreateNotes, GetNotes, UpdateNotes, DeleteNotes };
