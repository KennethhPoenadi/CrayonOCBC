export type Article = {
  id: number;
  documentId: string;
  title: string;
  description: string;
  slug: string;
  cover?: {
    url: string;
  };
  author?: {
    id: number;
    name: string;
  };
  category?: {
    id: number;
    name: string;
  };
};