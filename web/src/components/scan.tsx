"use client";

import { Button } from "@/components/ui/button";
import { Scan } from "lucide-react";
import { Scanner } from "@yudiel/react-qr-scanner";
import Modal from "./modal";

export default function QrScanner() {
  return (
    <Modal
      trigger={
        <Button variant="secondary">
          <Scan />
          Scan
        </Button>
      }
    >
      <Scanner
        onScan={(results) => (window.location.href = results[0]?.rawValue)}
      />
    </Modal>
  );
}
