import BackButton from "../components/general/buttons/backButton";
import BasicHeading from "../components/general/headings/basicHeading";

export default function Music() {
  return (
    <div>
      <BasicHeading text="Music Page" />
      <BackButton linkTo="/" />
    </div>
  );
}
