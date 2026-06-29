import { useParams } from "react-router-dom";
import useFetch from "../../hooks/useFetch";
import type { Article } from "../../types/article";

export default function ArticleDetailPage() {
  const { id } = useParams();

  const { data, loading, error } = useFetch<{ data: Article }>(
    `articles/${id}?populate=*`
  );

  if (loading) return <main>Loading...</main>;
  if (error) return <main>Error: {error}</main>;

  const article = data?.data;

  if (!article) return <main>Article not found</main>;

  return (
    <main>
      {article.cover?.url && (
        <img
          src={`http://localhost:1337${article.cover.url}`}
          alt={article.title}
          width={300}
        />
      )}

      <h1>{article.title}</h1>
      <p>{article.description}</p>

      <p>Slug: {article.slug}</p>
      <p>Author: {article.author?.name}</p>
      <p>Category: {article.category?.name}</p>
    </main>
  );
}