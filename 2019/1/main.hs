main = interact processInput

processInput input = show $ sum [perLine line | line <- lines input]

perLine :: String -> Integer
perLine line = fuel
  where fuel = fuelRequiredSimple val
        val = read line


fuelRequiredSimple :: Integer -> Integer
fuelRequiredSimple mass
  | amt < 0 = 0
  | otherwise = amt
  where amt = mass `div` 3 - 2
