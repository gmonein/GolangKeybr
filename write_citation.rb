File.read('citations.txt').split("\r\n\r\n").each.with_index { |s, i| File.write('citations/' + i.to_s, s.gsub('’', '\'').gsub("\r\n", "\n").gsub('—', '-')) }
