import Button from "@mui/material/Button";
import { Link } from "react-router-dom";

export default function BackButton({ linkTo }) {
  return (
    <>
      <Button component={Link} to={linkTo} variant="contained">
        Back
      </Button>
    </>
  );
}
