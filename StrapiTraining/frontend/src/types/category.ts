import type { Article } from "./article"


export type Category = {
    id: number;
    documentId: string;
    name: string;
    slug: string;
    articles: Article[];
    description: string;
}