import { Bars } from "react-loader-spinner";

export default function Loading() {
  return (
    <div className="fixed inset-0 bg-black bg-opacity-100 z-50 flex justify-center items-center">
      <Bars
        height="70"
        width="70"
        color="#e2714a"
        ariaLabel="bars-loading"
        visible={true}
      />
    </div>
  );
}
