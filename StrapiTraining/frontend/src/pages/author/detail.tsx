import type { Author } from "../../types/author"
import { useParams } from "react-router-dom";
import useFetch from "../../hooks/useFetch";

export default function AuthorDetailPage() {

    const { id } = useParams();

    const {data, loading, error} = useFetch<{data: Author}>(
        `authors/${id}?populate=*`
    );

    if (loading) return <main>Loading...</main>
    if (error) return <main>Error: {error} </main>

    const author = data?.data;

    if (!author) return <main>Author not found</main>;

    return (
        <main>
            {author.avatar?.url && (
            <img
                src={`http://localhost:1337${author.avatar.url}`}
                alt={author.name}
                width={200}
            />
            )}

            <h2>{author.name}</h2>
            <p>{author.email}</p>
        </main>
    );
}