package rongta

func (p *Printer) Println(text string) error {
	err := p.driver.WriteStringToBuffer(text + "\n")
	if err != nil {
		return err
	}

	// TODO: Calculate the number of lines to feed
	err = p.driver.PrintAndFeedNLines(10)
	if err != nil {
		return err
	}

	return nil
}
