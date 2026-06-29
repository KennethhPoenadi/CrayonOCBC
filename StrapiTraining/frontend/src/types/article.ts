import type { Author } from "./author"
import type { Category } from "./category";

export type Article = {
  id: number;
  documentId: string;
  title: string;
  description: string;
  slug: string;
  cover?: {
    url: string;
  };
  author: Author;
  category: Category;
};