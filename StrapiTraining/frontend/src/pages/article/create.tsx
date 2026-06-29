import { useState } from "react";
import { Button, Card, Form, Input, Select, message } from "antd";
import { api } from "../../lib/api";
import useFetch from "../../hooks/useFetch";
import type { Author } from "../../types/author";
import type { Category } from "../../types/category";

type ArticleFormValues = {
  title: string;
  description: string;
  slug: string;
  author?: string;
  category?: string;
};

export default function ArticleCreatePage() {
  const [submitting, setSubmitting] = useState(false);

  const { data: authors } = useFetch<{ data: Author[] }>("/authors");
  const { data: categories } = useFetch<{ data: Category[] }>("/categories");

  async function handleSubmit(values: ArticleFormValues) {
    try {
      setSubmitting(true);

      await api.post("/articles", {
        data: {
          title: values.title,
          description: values.description,
          slug: values.slug,
          author: values.author,
          category: values.category,
        },
      });

      message.success("Article created successfully");
    } catch (err) {
      message.error("Failed to create article");
      console.error(err);
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <main style={{ padding: 24 }}>
      <Card title="Create Article" style={{ maxWidth: 720 }}>
        <Form layout="vertical" onFinish={handleSubmit}>
          <Form.Item
            label="Title"
            name="title"
            rules={[{ required: true, message: "Title wajib diisi" }]}
          >
            <Input placeholder="Beautiful picture" />
          </Form.Item>

          <Form.Item
            label="Description"
            name="description"
            rules={[{ required: true, message: "Description wajib diisi" }]}
          >
            <Input.TextArea rows={4} placeholder="Description..." />
          </Form.Item>

          <Form.Item
            label="Slug"
            name="slug"
            rules={[{ required: true, message: "Slug wajib diisi" }]}
          >
            <Input placeholder="beautiful-picture" />
          </Form.Item>

          <Form.Item label="Author" name="author">
            <Select
              placeholder="Select author"
              options={authors?.data.map((author) => ({
                label: author.name,
                value: author.documentId,
              }))}
            />
          </Form.Item>

          <Form.Item label="Category" name="category">
            <Select
              placeholder="Select category"
              options={categories?.data.map((category) => ({
                label: category.name,
                value: category.documentId,
              }))}
            />
          </Form.Item>

          <Button type="primary" htmlType="submit" loading={submitting}>
            Create Article
          </Button>
        </Form>
      </Card>
    </main>
  );
}