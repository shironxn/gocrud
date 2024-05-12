"use server";

import { NoteCreate, NoteQuery, NoteUpdate } from "@/lib/schema/note";
import { revalidatePath } from "next/cache";
import { cookies, headers } from "next/headers";
import { notFound } from "next/navigation";

const BASE_API_URL = process.env.NEXT_PUBLIC_API_URL;

const CreateNotes = async (data: NoteCreate) => {
  const res = await fetch(`${BASE_API_URL}/notes`, {
    method: "POST",
    body: JSON.stringify(data),
    headers: {
      "Content-Type": "application/json",
      Cookie: cookies().toString(),
    },
  });

  if (!res.ok) {
    const result = await res.json();
    return result.error;
  }

  revalidatePath("/");
};

const GetNotes = async (query?: NoteQuery) => {
  const res = await fetch(
    `${BASE_API_URL}/notes?user_id=${query?.user_id || 0}&visibility=${
      query?.visibility || ""
    }&title=${query?.search || ""}&page=${query?.page || 1}&limit=6&order=desc`,
    { headers: headers() }
  );

  const result = await res.json();
  if (!res.ok) {
    return { error: result.error };
  }

  return result;
};

const GetNotesByID = async (id: string) => {
  const res = await fetch(`${BASE_API_URL}/notes/${id}`, {
    headers: headers(),
  });

  if (res.status === 404) {
    notFound();
  }

  const result = await res.json();
  if (!res.ok) {
    return { error: result?.error };
  }

  return result;
};

const UpdateNotes = async (data: NoteUpdate, id: number) => {
  const res = await fetch(`${BASE_API_URL}/notes/${id}`, {
    method: "PUT",
    body: JSON.stringify(data),
    headers: {
      "Content-Type": "application/json",
      Cookie: cookies().toString(),
    },
  });

  if (!res.ok) {
    const result = await res.json();
    return result.error;
  }

  revalidatePath("/");
};

const DeleteNotes = async (id: string) => {
  const res = await fetch(`${BASE_API_URL}/notes/${id}`, {
    method: "DELETE",
    headers: headers(),
  });

  if (!res.ok) {
    const result = await res.json();
    return result.error;
  }
};

export { CreateNotes, GetNotes, GetNotesByID, UpdateNotes, DeleteNotes };
