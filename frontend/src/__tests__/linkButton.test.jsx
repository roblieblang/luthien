import "@testing-library/jest-dom";
import { fireEvent, render, screen } from "@testing-library/react";
import React from "react";
import { BrowserRouter } from "react-router-dom";
import LinkButton from "../components/general/buttons/linkButton";

describe("LinkButton component", () => {
  test("renders an external link correctly", () => {
    render(<LinkButton to="http://external.com" text="External Link" />);
    const linkElement = screen.getByText(/External Link/i);
    expect(linkElement).toBeInTheDocument();
    expect(linkElement).toHaveAttribute("href", "http://external.com");
    expect(linkElement).toHaveAttribute("target", "_blank");
    expect(linkElement).toHaveAttribute("rel", "noopener noreferrer");
  });

  test("renders an internal link correctly", () => {
    render(
      <BrowserRouter>
        <LinkButton to="/internal" text="Internal Link" />
      </BrowserRouter>
    );
    const linkElement = screen.getByText(/Internal Link/i);
    expect(linkElement).toBeInTheDocument();
    expect(linkElement.closest("a")).toHaveAttribute("href", "/internal");
  });

  test('renders a button when no "to" prop is provided', () => {
    render(<LinkButton text="Button" />);
    const buttonElement = screen.getByText(/Button/i);
    expect(buttonElement).toBeInTheDocument();
    expect(buttonElement.tagName).toBe("BUTTON");
  });

  test("calls onClick handler when clicked", () => {
    const handleClick = jest.fn();
    render(<LinkButton text="Button" onClick={handleClick} />);
    const buttonElement = screen.getByText(/Button/i);
    fireEvent.click(buttonElement);
    expect(handleClick).toHaveBeenCalledTimes(1);
  });
});
